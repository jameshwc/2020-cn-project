#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <stdlib.h>
#include "config.h"

int main(){
    struct sockaddr_in server_addr;
    struct sockaddr_storage server_storage;
    char buffer[1024];
    socklen_t addr_size;
    int welcome_socket, new_socket, client_len, n_bytes;
    
    welcome_socket = socket(PF_INET, SOCK_STREAM, 0);
    
    server_addr.sin_family = AF_INET;
    server_addr.sin_port = htons(PORT_NUM);
    server_addr.sin_addr.s_addr = inet_addr(SERVER_ADDR);
    memset(server_addr.sin_zero, '\0', sizeof(server_addr.sin_zero));
    
    bind(welcome_socket, (struct sockaddr *) &server_addr, sizeof(server_addr));

    if(listen(welcome_socket, MAX_CONN)==0) printf("Listening\n");
    else printf("Error\n");

    addr_size = sizeof(server_storage);
    while(1){
        new_socket = accept(welcome_socket, (struct socketaddr *) &server_storage, &addr_size);
        if(!fork()){
            n_bytes = 1;
            while(n_bytes!=0){
                n_bytes = recv(new_socket,buffer,1024,0);
    
                for (int i = 0; i < n_bytes-1; i++){
                    buffer[i] = toupper(buffer[i]);
                }
                printf("%s\n", buffer);
                send(new_socket,buffer, n_bytes, 0);
            }
            close(new_socket);
            exit(0);
        }
        else {
            close(new_socket);
        }
  }

  return 0;
}