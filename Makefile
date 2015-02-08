all: jsx cleanup embed build

check-react:
	@type jsx >/dev/null 2>&1 || (echo "Please install React.js toolkit:\nnpm install -g react-tools\n" && exit 1)

jsx: check-react
	jsx --extension jsx server/static/ server/static/
	rm -rf server/static/.module-cache

watch:
	jsx --watch --extension jsx server/static/ server/static/

cleanup:
	find . -name '.DS_Store' | xargs rm -f

embed:
	rice -i ./server embed-go

build:
	go build
