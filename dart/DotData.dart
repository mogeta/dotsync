class DotData {
  int width;
  int height;
  List listR;
  List listG;
  List listB;
  
  DotData(int w,int h){
    width  = w;
    height = h;
    listR = new List(w*h);
    listG = new List(w*h);
    listB = new List(w*h);
  }
}