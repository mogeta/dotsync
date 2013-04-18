import 'package:js/js.dart' as js;
import 'package:appengine_channel/appengine_channel.dart' as ch;
import 'dart:html';

void main() {
  
  String token;
  String meID;
  String gameID;
  
  //Exception: Cannot get JavaScript context out of scope.
  //token = js.context.token;
  
  js.scoped(() {
     token  = js.context.token;
     meID   = js.context.me;
     gameID = js.context.game_key;
   });
  
  openChannel(token);
  setButtonListner(meID,gameID);
  
}

void openChannel(String token) {
  Element element = document.query("#area");
  ch.Channel channel = new ch.Channel(token);
  ch.Socket   socket = channel.open()
    ..onOpen    = (() => print("open"))
    ..onClose   = (() => print("close"))
    ..onMessage = ((m){ 
      element.innerHtml = "${element.innerHtml}${m}<br />";
        print("${m}");
      })
    ..onError   = ((code, desc) => print("error: $code $desc"));
}

void setButtonListner(String meID,String gameID){
  var btn;
  btn = document.query("#button_id");

  btn.onClick.listen((e){
    var httpRequest;
    httpRequest = new HttpRequest();
    httpRequest.open('POST', '/receive?g=${meID}${gameID}', async:true);
    httpRequest.send();
  });
}