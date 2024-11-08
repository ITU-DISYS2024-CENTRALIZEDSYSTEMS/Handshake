# Handshake

## Questions
### Question A
_What are packages in your implementation? What data structure do you use to transmit data and meta-data?_

We don't handle any packages in our implementation. We only use strings to transmit data and meta-data with the use of the `net` package. 

### Question B 
_Does your implementation use threads or processes? Why is it not realistic to use threads?_

We make use of goroutines in our implementation, but we have structured the goroutines as a 'client' & 'server' relationship. It is not realistic to use threads in the real world, because TCP/UDP is a internet protocol between a server & client - We only use threads in our example to simulate the client and server relationship.

### Question C 
_In case the network changes the order in which messages are delivered, how would you handle message re-ordering?_

As we only handle a 3-way handshake and only send a single package at a time, the order of the packages will always be the same. As for lost packages, the program crashes if any of the expected packages does not arrive.

Our program also does not support sending messages (for example data), but we would handle the arrival of packets in the wrong order by checking the sqeuence numbers. Each packet have a sequence number, so we would be able to buffer the packets when receiving and reorder when every packet is received. 

### Question D 
_In case messages can be delayed or lost, how does your implementation handle message loss?_

Our current implementation does not handle this scenario. But it would be achieved by the parties validating the sequence number whenever they receive a packet. If the sequence number does not match that packet (and packets sent after) would need to be resent.

### Question E 
_Why is the 3-way handshake important?_

The handshake is important because it ensures that both parties (in this case the server and the client) are ready to send and receive data. But also that the sequence numbers are correct, so that the messages can be acknowledged


