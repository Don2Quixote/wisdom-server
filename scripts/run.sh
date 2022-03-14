source scripts/colors.sh

check() {
    if test $? -ne 0
    then
    	printf "$*"
        exit 1
    fi

}

printf "${YELLOW}building image...${NC}\n"
sudo docker build -f build/Dockerfile -t wisdom-server .
check "${RED}failed to build image${NC}\n"

printf "${YELLOW}launching container...${NC}\n"
sudo docker run --name wisdom-server --rm \
    -e PORT=4444 \
    -e COMPLEXITY_FACTOR=0.25 \
    -e MAX_COMPLEXITY=16 \
    -e COMPLEXITY_DURATION_SECONDS=44 \
    -p 4444:4444 wisdom-server
check "${RED}failed to launch container${NC}\n"

