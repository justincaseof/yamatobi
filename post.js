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
   /*request.open("POST","http://192.168.178.7/YamahaRemoteControl/ctrl");		// CORS :-(
   request.setRequestHeader("Content-Length",""+data.length);
   request.setRequestHeader("Content-Type","text/xml");
   request.setRequestHeader("Origin", 'FOO');
   request.send(data);
   */
   console.log("## >>>> Preset: " + preset);
   request.open("GET","http://127.0.0.1:9000/preset/" + preset);
   request.send();
   console.log("## <<<<");
}