#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <fcntl.h>
#include <signal.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/stat.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include "config.h"

void convert_ip_addr(in_addr_t addr, char* addr_str) {
    inet_ntop(AF_INET, &(addr), addr_str, INET_ADDRSTRLEN);
}

void handle_socket(int socket_fd) {
    static char buffer[BUF_SIZE+1];

    int n = read(socket_fd,buffer,BUF_SIZE); 
    if (n <= 0) {
        exit(3);
    }

    if (n < BUF_SIZE) buffer[n] = 0;

    for (int i = 0; i < n; i++) 
        if (buffer[i]=='\r'||buffer[i]=='\n')
            buffer[i] = 0;
    
    if (strncmp(buffer,"GET ",4)&&strncmp(buffer,"get ",4))
        exit(3);
    
    for(int i = 4; i < BUF_SIZE; i++) {
        if(buffer[i] == ' ') {
            buffer[i] = 0;
            break;
        }
    }

    for (int j = 0; j < BUF_SIZE; j++){
        if(buffer[j] == 0) break;
        if (buffer[j]=='.' && buffer[j+1]=='.') exit(3);
    }

    if (!strncmp(&buffer[0],"GET /\0",6)||!strncmp(&buffer[0],"get /\0",6) )
        strcpy(buffer,"GET /index.html\0");

    int file_fd = open(&buffer[5],O_RDONLY);
    if(file_fd==-1) write(socket_fd, "Failed to open file", 19);

    struct stat buf;
    fstat(file_fd, &buf);

    sprintf(buffer,"HTTP/1.0 200 OK\r\nContent-Type: text/html\r\nContent-Length: %d\r\n\r\n", buf.st_size);
    write(socket_fd,buffer,strlen(buffer));

    while ((n=read(file_fd, buffer, BUF_SIZE))>0) {
        write(socket_fd,buffer, n);
    }

    exit(1);
}

int main() {
    int listen_fd, socket_fd;
    struct sockaddr_in cli_addr;
    struct sockaddr_in serv_addr;
    signal(SIGCLD, SIG_IGN);
    listen_fd = socket(AF_INET, SOCK_STREAM,0);
    if (listen_fd < 0) exit(3);

    serv_addr.sin_addr.s_addr = htonl(INADDR_ANY);
    serv_addr.sin_port = htons(PORT_NUM);
    serv_addr.sin_family = AF_INET;

    char ip_addr[INET_ADDRSTRLEN];
    convert_ip_addr(serv_addr.sin_addr.s_addr, ip_addr);
    int one = 1;
    setsockopt(listen_fd, SOL_SOCKET, SO_REUSEADDR, &one, sizeof(one));  
    while (bind(listen_fd, (struct sockaddr *)&serv_addr,sizeof(serv_addr))<0) {
        printf("trying to bind server address %s\n", ip_addr);
        sleep(1);
    };
    printf("binding success!\n");

    if (listen(listen_fd, 64) < 0) exit(5);

    while(1) {
        pid_t pid;
        int length = sizeof(cli_addr);
        socket_fd = accept(listen_fd, (struct sockaddr *)&cli_addr, &length);
        if (socket_fd < 0) exit(6);
        
        if ((pid = fork()) < 0) {
            exit(7);
        } else {
            if (pid == 0) {  // child process
                close(listen_fd);
                handle_socket(socket_fd);
            } else { // parent process
                char ip_addr[INET_ADDRSTRLEN];
                convert_ip_addr(cli_addr.sin_addr.s_addr, ip_addr);
                printf("GET connection from %s\n",ip_addr);
                close(socket_fd);
            }
        }
    }
}

