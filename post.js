/*
var form = document.querySelector('form');
form.addEventListener('submit', function (event) {
   event.preventDefault();
   var data = new FormData(form);
   check(data);
});
*/

function send (preset){
   var request = new XMLHttpRequest();
   request.addEventListener('load', function(event) {
      if (request.status >= 200 && request.status < 300) {
         console.log(request.responseText);
      } else {
         console.warn(request.statusText, request.responseText);
      }
   });
   console.log("## >>>> Preset: " + preset);
   request.open("GET","http://127.0.0.1:9000/preset/" + preset);
   request.send();
   console.log("## <<<<");
}

function shutdown (){
   var request = new XMLHttpRequest();
   request.addEventListener('load', function(event) {
      if (request.status >= 200 && request.status < 300) {
         console.log(request.responseText);
      } else {
         console.warn(request.statusText, request.responseText);
      }
   });
   request.open("GET","http://127.0.0.1:9000/exit");
   request.send();
}