resource "aws_vpc" "pavement" {
    cidr_block = "10.0.0.0/16"
    enable_dns_hostnames = true
    enable_dns_support = true
}

resource "aws_subnet" "pavement" {
    cidr_block = "${cidrsubnet(aws_vpc.pavement.cidr_block, 3, 1)}"  
    vpc_id = "${aws_vpc.pavement.id}"
    availability_zone = "us-east-1a"
}

resource "aws_security_group" "pavement" {
    name = "allow-all-sg"
    vpc_id = "${aws_vpc.pavement.id}"

    ingress {
        cidr_blocks = [
            "0.0.0.0/0"
        ]

        from_port = 22
        to_port = 22
        protocol = "tcp"       

    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}