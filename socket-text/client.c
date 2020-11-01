#include <stdio.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include "config.h"

int main(){
  int client_socket, n_bytes;
  char buffer[1024];
  struct sockaddr_in server_addr;
  socklen_t addr_size;

  client_socket = socket(PF_INET, SOCK_STREAM, 0);
  server_addr.sin_family = AF_INET;
  server_addr.sin_port = htons(PORT_NUM);
  server_addr.sin_addr.s_addr = inet_addr(SERVER_ADDR);
  memset(server_addr.sin_zero, '\0', sizeof(server_addr.sin_zero));

  addr_size = sizeof(server_addr);
  connect(client_socket, (struct sockaddr *) &server_addr, addr_size);

  while(1){
    printf("Type a sentence to send to server:\n");
    fgets(buffer,1024,stdin);
    printf("You typed: %s",buffer);

    n_bytes = strlen(buffer) + 1;

    send(client_socket,buffer,n_bytes,0);

    recv(client_socket, buffer, 1024, 0);

    printf("Received from server: %s\n\n",buffer);   
  }

  return 0;
}