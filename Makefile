SHELL			:= /bin/bash

build:
	for example in simple ec2 custom-logs; do make -C $${example} build; done

setup:
	for example in simple ec2 custom-logs; do make -C $${example} setup; done

clean-setup:
	for example in simple ec2 custom-logs; do make -C $${example} clean-setup; done

all: setup build clean-setup