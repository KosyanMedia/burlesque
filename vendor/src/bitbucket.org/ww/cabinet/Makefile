include $(GOROOT)/src/Make.inc

TARG=bitbucket.org/ww/cabinet
CGOFILES=cabinet.go		

include $(GOROOT)/src/Make.pkg

format:
	for x in *.go; do gofmt $${x} > $${x}.fmt && mv $${x}.fmt $${x}; done

docs:
	gomake clean
	godoc ${TARG} > README.txt
