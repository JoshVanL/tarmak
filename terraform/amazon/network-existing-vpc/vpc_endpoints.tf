resource "aws_vpc_endpoint" "s3" {
  vpc_id          = "${var.vpc_id}"
  service_name    = "com.amazonaws.${var.region}.s3"
  route_table_ids = ["${aws_route_table.private.*.id}"]
}
