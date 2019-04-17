
running_image: image
	(docker stop donaldguyartsrv && docker rm -fv donaldguyartsrv) || :
	docker run -d -p 4444:4444 --name donaldguyartsrv recruitment/donaldguy

image:
	docker build -t recruitment/donaldguy .
	
test:
ifeq ($(OS),Windows_NT)
	testbins/ascii_windows_x64.exe
else
	testbins/ascii_$(shell uname -s | tr '[:upper:]' '[:lower:]')_x64
endif

logs:
	docker logs -f donaldguyartsrv

.PHONY: running_image image test logs
