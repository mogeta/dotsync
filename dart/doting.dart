import 'package:js/js.dart' as js;
import 'package:appengine_channel/appengine_channel.dart' as ch;
import 'package:color_picker/color_picker.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';
import 'Channel.dart';

class Doting { 
  CanvasRenderingContext2D ctx;
  Channel ch;
  var clientRect;
  
  int nowR=0,nowG=0,nowB=0;
  
  Doting() {
    CanvasElement canvas = document.query("#canvas");

    ch = new Channel();
    ch.receive  = receive;
    ctx = canvas.getContext("2d");  
    
    var colorPicker = new ColorPicker(256,showInfoBox: true);
    document.query("#picker").nodes.add(colorPicker.element);
    
    
    colorPicker.colorChangeListener = (ColorValue color, num hue, num saturation, num brightness) {
        // Process new color value here ...
      nowR=color.r;
      nowG=color.g;
      nowB=color.b;
    };
    
    clientRect = ctx.canvas.getBoundingClientRect();
    ctx.canvas.onClick.listen(canvasClickEvent);
  }
  
  void receive(m){
    print("test");
    print(m);
    Map data = parse(m);
    drawPixel(data["x"],data["y"],data["r"],data["g"],data["b"]);
  }
  
  void canvasClickEvent(e){
    var x = ((e.clientX - clientRect.left)/32).floor();
    var y = ((e.clientY - clientRect.top)/32).floor();
    drawPixel(x,y,nowR,nowG,nowB);
    
    var mapData = new Map();
    mapData["x"]  = x;
    mapData["y"]  = y;
    mapData["r"]  = nowR;
    mapData["g"]  = nowG;
    mapData["b"]  = nowB;
    ch.sendMessage(stringify(mapData));
  }
  
  void drawPixel(int x,int y,int setR,int setG,int setB){
    const int LIM=10;
    int cnt=0;
    double r=255.0,g=255.0,b=255.0;
    double  dr=(r-setR)/(-10);
    double  dg=(g-setG)/(-10);
    double  db=(b-setB)/(-10);
    
    new Timer.periodic(new Duration(milliseconds: 100),
        (Timer t){
          cnt++;
          r+=dr;
          g+=dg;
          b+=db;
          //print('rgb(${r.round()}, ${g.round()}, ${b.round()})');
          ctx.fillStyle = 'rgb(${r.round()}, ${g.round()}, ${b.round()})';
          ctx.fillRect(x*32, y*32, 32, 32);
          if(cnt>=LIM){
            t.cancel();
          }
        });
    
    //create sendData

    //runAsync(print(stringify(mapData)));
  }

}

void main() {  
  new Doting();
}  