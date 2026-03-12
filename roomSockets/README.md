Okkay I was hovering around excalidraw's code repo and found so much interesting things from their collaboration feature 

1/n -- event driven socket architecture
- the server never sees plaintext, the initialization vector and roomKey pattern is AES-GCM encryption 

some things to really get from this
> the server doesn't need to read data to route it
> e2e encryption at the socket layer
> iv is sent alongside data so the receiver can decrypt without a shared session

And this is also got divided into layers 

JUSTT WAOOWWWW