# the command macro
GOBUILD := go build -ldflags '-s'
SCP := scp -P 10022
ALIYUN := sunchao@112.126.66.162:/root/workspace/hf

# the elf program names
SERVER := herefriend
CRAW := spider
TRIM := trimimage
MODIFY := modify
SENDGIFT := sendgift

SERVER_GOFILES := \
	$(wildcard ./common/*.go) \
	$(wildcard ./config/*.go) \
	$(wildcard ./lib/*.go) \
	$(wildcard ./lib/push/*.go) \
   	$(wildcard ./server/handlers/*.go) \
	$(wildcard ./server/routes/*.go) \
	$(wildcard ./server/cms/*.go) \
	$(wildcard ./server/*.go) 

CRAW_GOFILES := \
	$(wildcard ./common/*.go) \
	$(wildcard ./config/*.go) \
	$(wildcard ./lib/*.go) \
	$(wildcard ./lib/push/*.go) \
	$(wildcard ./crawler/dbtables/*.go) \
	$(wildcard ./crawler/idsearch/*.go) \
	$(wildcard ./crawler/image/*.go) \
	$(wildcard ./crawler/page/*.go) \
	$(wildcard ./crawler/request/*.go) \
	$(wildcard ./crawler/page3g/*.go) \
	$(wildcard ./crawler/pageweibo/*.go) \
	$(wildcard ./crawler/pagezhenqing/*.go) \
	$(wildcard ./crawler/filter/*.go) \
	$(wildcard ./crawler/*.go)

TRIM_GOFILES := \
	$(wildcard ./common/*.go) \
	$(wildcard ./config/*.go) \
	$(wildcard ./crawler/image/*.go) \
	$(wildcard ./lib/*.go) \
	$(wildcard ./tools/imageproc/*.go)

MODIFY_GOFILES := \
	$(wildcard ./common/*.go) \
	$(wildcard ./config/*.go) \
	$(wildcard ./lib/*.go) \
	$(wildcard ./tools/modifier/*.go)

SENDGIFT_GOFILES := \
	$(wildcard ./common/*.go) \
	$(wildcard ./config/*.go) \
	$(wildcard ./lib/*.go) \
	$(wildcard ./tools/giftsender/*.go)

server: $(SERVER)
trim: $(TRIM)
craw: $(CRAW)
modify: $(MODIFY)
gift: $(SENDGIFT)

$(SERVER): $(SERVER_GOFILES)
	@$(GOBUILD) -o ./$@ herefriend/server
	@echo "finish"

x64: $(SERVER_GOFILES)
	@GOOS=linux $(GOBUILD) -o ./$(SERVER) herefriend/server
	@echo "finish"

$(CRAW): $(CRAW_GOFILES)
	@$(GOBUILD) -o ./$@ herefriend/crawler
	@echo "finish"

craw.x64: $(CRAW_GOFILES)
	@GOOS=linux $(GOBUILD) -o ./$(CRAW) herefriend/crawler
	@echo "finish"

$(TRIM): $(TRIM_GOFILES)
	@$(GOBUILD) -o ./$@ herefriend/tools/imageproc
	@echo "finish"

$(MODIFY): $(MODIFY_GOFILES)
	@$(GOBUILD) -o ./$@ herefriend/tools/modifier
	@echo "finish"

$(SENDGIFT): $(SENDGIFT_GOFILES)
	@$(GOBUILD) -o ./$@ herefriend/tools/giftsender
	@echo "finish"

gift.x64: $(SENDGIFT_GOFILES)
	@GOOS=linux $(GOBUILD) -o ./$(SENDGIFT) herefriend/tools/giftsender
	@echo "finish"
	
all: $(SERVER) $(CRAWLER)

clean:
	@rm -rf ./$(SERVER) ./$(SERVER).pid ./$(CRAW) ./$(TRIM) ./$(MODIFY) ./$(SENDGIFT) ./log
	@echo "finish"

cp:
	@tar czf $(SERVER).gz $(SERVER)
	@$(SCP) $(SERVER).gz $(ALIYUN)
	@rm -f $(SERVER).gz $(SERVER)
	@echo "finish"

cpcraw:
	@tar czf $(CRAW).gz $(CRAW)
	@$(SCP) $(CRAW).gz $(ALIYUN)
	@rm -f $(CRAW).gz $(CRAW)
	@echo "finish"

cptrim:
	@tar czf $(TRIM).gz $(TRIM)
	@$(SCP) $(TRIM).gz $(ALIYUN)
	@rm -f $(TRIM).gz $(TRIM)
	@echo "finish"

cpmodify:
	@tar czf $(MODIFY).gz $(MODIFY)
	@$(SCP) $(MODIFY).gz $(ALIYUN)
	@rm -f $(MODIFY).gz $(MODIFY)
	@echo "finish"

cpgift:
	@tar czf $(SENDGIFT).gz $(SENDGIFT)
	@$(SCP) $(SENDGIFT).gz $(ALIYUN)
	@rm -f $(SENDGIFT).gz $(SENDGIFT)
	@echo "finish"
