SHELL		 := /bin/bash

ECHO			= echo -e
C_STD			= $(shell echo -e "\033[0m")
C_RED			= $(shell echo -e "\033[91m")
C_GREEN		= $(shell echo -e "\033[92m")
C_YELLOW	= $(shell echo -e "\033[93m")
P					= $(shell echo -e "\033[92m> \033[0m")
OK				= $(shell echo -e "\033[92m[ OK  ]\033[0m")
ERROR			= $(shell echo -e "\033[91m[ERROR] \033[0m")
PASS			= $(shell echo -e "\033[92m[PASS ]\033[0m")
FAIL			= $(shell echo -e "\033[91m[FAIL ] \033[0m")

default:
	@export GOPATH=$$(mktemp -d); \
	$(ECHO) "$(C_GREEN)Building temporary Go modules cache $(C_STD)$${GOPATH}/pkg/mod"; \
	go clean -modcache; \
	for example in 0*; do \
		cd $${example}; \
		go mod download; \
		$(ECHO) "Building $${example} ..."; \
		$(RM) terractl_test; \
		go build -o terractl_test .; \
		if [[ -e terractl_test ]]; then \
			$(ECHO) "$(OK) $(C_GREEN)$${example}$(C_STD)"; \
		else \
			$(ECHO) "$(ERROR) $(C_RED)$${example}$(C_STD)"; \
		fi; \
		$(RM) terractl_test; \
		cd ..; \
	done; \
	$(ECHO) "$(C_GREEN)Cleaning temporary Go modules cache directory $(C_STD)$${GOPATH}/pkg/mod:$(C_STD)"; \
	go clean -modcache; \
	$(RM) -r $${GOPATH}

tidy:
	for example in 0*; do cd $${example}; go mod tidy; cd ..; done

build:
	for example in 0*; do [[ ! -e $${example}/Makefile ]] || make -C $${example}; done

setup:
	[[ -e ../terranova ]]
	for example in 0* do [[ ! -e $${example}/Makefile ]] || make -C $${example} setup; done

clean-setup:
	[[ -e ../terranova ]]
	for example in 0* do [[ ! -e $${example}/Makefile ]] || make -C $${example} clean-setup; done

all: setup build clean-setup