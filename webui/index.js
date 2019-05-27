//https://www.compose.com/articles/redis-pubsub-node-and-socket-io/

var app = require('express')();
var http = require('http').Server(app);
var io = require('socket.io')(http);
var redis = require('redis');
//var url = config.get('redis.url');

//incase of docker need to specify redis server separately
var client1 = redis.createClient(6379, "172.25.16.126");

app.get('/', function(req, res) {
   res.sendfile('index.html');
});

io.on('connection', onConnection);

http.listen(3000, function() {
   console.log('listening on *:3000');
});

function onConnection(socket){
  console.log('A user connected');

  setTimeout(function(){
    socket.send("Message sent after 4000ms")
  }, 1000);

  client1.on('message', function(chan, msg) {
    console.log(msg);
    socket.emit('redisPublishEvent', {description : msg});
  });

  var value = 10;
  //Send a message when
  setInterval(function() {
     //Sending an object when emmiting an event

     socket.emit('testCustomEvent', { description: 'A custom event named testerEvent!', data1 : value++, data2 :value++, data3 : value+2, data4 : value+3});
  }, 2000);

  socket.on('disconnect', function(){
    console.log('A user disconnected')
  })
}

client1.subscribe('hashChannel');
