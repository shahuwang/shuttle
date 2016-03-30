server 端的逻辑：
1. 接收client端建立tunnel的请求, 建立server_hub
2. 从tunnel读取出linkid和data (tunnel的read共两次，第一次读取头部，
   包括linkid和数据长度，第二步读取剩余所有数据)
3. 
