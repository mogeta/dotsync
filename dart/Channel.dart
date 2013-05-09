library Channel;

import 'package:js/js.dart' as js;
import 'package:appengine_channel/appengine_channel.dart' as ch;
import 'dart:html';

typedef void ReceiveMessage(m);

class Channel { 
  
  ReceiveMessage receive;
  
  String token;
  String meID;
  String gameID;
  
  //Exception: Cannot get JavaScript context out of scope.
  //token = js.context.token;
  
  Channel(){
    js.scoped(() {
      token  = js.context.token;
      meID   = js.context.me;
      gameID = js.context.game_key;
    });
    
    openChannel(token);
  }
  
  void openChannel(String token) {
    Element element = document.query("#area");
    ch.Channel channel = new ch.Channel(token);
    ch.Socket   socket = channel.open()
        ..onOpen    = (() => print("open"))
        ..onClose   = (() => print("close"))
        ..onMessage = ((m) => receive(m))
        ..onError   = ((code, desc) => print("error: $code $desc"));
  }
  
  void sendMessage(String msg){
    var httpRequest;
    httpRequest = new HttpRequest();
    httpRequest.open('POST', '/receive?p=${meID}${gameID}&m=${msg}', async:true);
    httpRequest.send();
  }
}