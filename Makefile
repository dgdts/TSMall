goods_service:
	hz update -idl idl/goods.proto --customize_package=./template/package.yaml   

all:
	make goods_service