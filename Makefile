generate:
	hz update -idl idl/$(service).proto 

ts_mall_service:
	make service=ts_mall_service generate

all:
	make ts_mall_service