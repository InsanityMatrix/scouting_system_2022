var canvas = document.getElementById('game');
var context = canvas.getContext('2d');
window.onload = function(){
    make_field();
}
function make_field() {
   field_image = new Image();
   field_image.src = "field.png";
   field_image.onload = function(){
     context.drawImage(field_image, 0,0);
   };
   
}

canvas.addEventListener("click", function(evt){
  console.log("Clicked");
  let mousePos = getMousePos(canvas, evt);
  //Show a Scoring Menu
  document.getElementById("overlay").style.display = "block";
});
function getMousePos(canvas, evt) {
  var rect = canvas.getBoundingClientRect();
  return {
    x: evt.clientX - rect.left,
    y: evt.clientY - rect.top
  };
}