package glue_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglue "github.com/hashicorp/terraform-provider-aws/internal/service/glue"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccGlueCatalogTableOptimizer_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var catalogTableOptimizer glue.GetTableOptimizerOutput

	resourceName := "aws_glue_catalog_table_optimizer.test"

	dbName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	tName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.GlueServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCatalogTableOptimizerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogTableOptimizerConfig_basic(dbName, tName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogTableOptimizerExists(ctx, resourceName, &catalogTableOptimizer),
					acctest.CheckResourceAttrAccountID(resourceName, names.AttrCatalogID),
					resource.TestCheckResourceAttr(resourceName, names.AttrDatabaseName, dbName),
					resource.TestCheckResourceAttr(resourceName, names.AttrTableName, tName),
					resource.TestCheckResourceAttr(resourceName, "type", "compaction"),
					// resource.TestCheckResourceAttr(resourceName, "configuration.role_arn", "arn:aws:iam::123456789012:role/example-role"),
					resource.TestCheckResourceAttr(resourceName, "configuration.enabled", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueCatalogTableOptimizer_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var catalogTableOptimizer glue.GetTableOptimizerOutput

	resourceName := "aws_glue_catalog_table_optimizer.test"

	dbName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	tName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.GlueServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckCatalogTableOptimizerDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccCatalogTableOptimizerConfig_basic(dbName, tName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCatalogTableOptimizerExists(ctx, resourceName, &catalogTableOptimizer),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfglue.ResourceCatalogTableOptimizer(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckCatalogTableOptimizerExists(ctx context.Context, resourceName string, catalogTableOptimizer *glue.GetTableOptimizerOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Glue Catalog Table Optimizer ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn(ctx)
		idParts := strings.Split(rs.Primary.ID, ":")
		if len(idParts) != 4 {
			return fmt.Errorf("unexpected format of ID (%q), expected catalog_id:database_name:table_name:type", rs.Primary.ID)
		}
		catalogID, databaseName, tableName, optimizerType := idParts[0], idParts[1], idParts[2], idParts[3]

		resp, err := conn.GetTableOptimizerWithContext(ctx, &glue.GetTableOptimizerInput{
			CatalogId:    aws.String(catalogID),
			DatabaseName: aws.String(databaseName),
			TableName:    aws.String(tableName),
			Type:         aws.String(optimizerType),
		})

		if err != nil {
			return fmt.Errorf("error getting Glue Catalog Table Optimizer (%s): %w", rs.Primary.ID, err)
		}

		if resp == nil || resp.TableOptimizer == nil {
			return fmt.Errorf("Glue Catalog Table Optimizer (%s) not found", rs.Primary.ID)
		}

		*catalogTableOptimizer = *resp

		return nil
	}
}

func testAccCheckCatalogTableOptimizerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_glue_catalog_table_optimizer" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn(ctx)
			idParts := strings.Split(rs.Primary.ID, ":")
			if len(idParts) != 4 {
				return fmt.Errorf("unexpected format of ID (%q), expected catalog_id:database_name:table_name:type", rs.Primary.ID)
			}
			catalogID, databaseName, tableName, optimizerType := idParts[0], idParts[1], idParts[2], idParts[3]

			_, err := conn.GetTableOptimizerWithContext(ctx, &glue.GetTableOptimizerInput{
				CatalogId:    aws.String(catalogID),
				DatabaseName: aws.String(databaseName),
				TableName:    aws.String(tableName),
				Type:         aws.String(optimizerType),
			})
			if err != nil {
				if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
					return nil
				}
				return fmt.Errorf("error getting Glue Catalog Table Optimizer (%s): %w", rs.Primary.ID, err)
			}

			return fmt.Errorf("Glue Catalog Table Optimizer %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCatalogTableOptimizerConfig_basic(dbName string, tName string) string {
	return fmt.Sprintf(`
data "aws_region" "current" {}
data "aws_caller_identity" "current" {}

resource "aws_iam_role" "test" {
  name               = %[1]q
  assume_role_policy = data.aws_iam_policy_document.assume.json
}

data "aws_iam_policy_document" "assume" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["glue.amazonaws.com"]
    }
  }
}

resource "aws_iam_role_policy" "GlueCompactionRoleAccess" {
  role = aws_iam_role.test.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": [
        "arn:aws:s3:::${aws_s3_bucket.bucket.bucket}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::${aws_s3_bucket.bucket.bucket}"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "glue:UpdateTable",
        "glue:GetTable"
      ],
      "Resource": [
        "arn:aws:glue:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:table/${aws_glue_catalog_database.test.name}/${aws_glue_catalog_table.test.name}",
        "arn:aws:glue:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:database/${aws_glue_catalog_database.test.name}",
        "arn:aws:glue:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:catalog"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:log-group:/aws-glue/iceberg-compaction/logs:*"
    }
  ]
}
EOF
}

resource "aws_glue_catalog_database" "test" {
  name = %[1]q
}

resource "aws_s3_bucket" "bucket" {
  bucket        = %[1]q
  force_destroy = true
}

resource "aws_glue_catalog_table" "test" {
  database_name = aws_glue_catalog_database.test.name
  name          = %[2]q
  table_type    = "EXTERNAL_TABLE"

  open_table_format_input {
    iceberg_input {
      metadata_operation = "CREATE"
      version            = 2
    }
  }

  storage_descriptor {
    location = "s3://${aws_s3_bucket.bucket.bucket}/files/"

    columns {
      name    = "my_column_1"
      type    = "int"
      comment = %[2]q
    }
  }
}

resource "aws_glue_catalog_table_optimizer" "test" {
  catalog_id     = data.aws_caller_identity.current.account_id
  database_name  = aws_glue_catalog_database.test.name
  table_name     = aws_glue_catalog_table.test.name
  configuration {
    role_arn = aws_iam_role.test.arn
    enabled  = true
  }
  type = "compaction"
}
`, dbName, tName)
}
