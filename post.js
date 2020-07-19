/*
var form = document.querySelector('form');
form.addEventListener('submit', function (event) {
   event.preventDefault();
   var data = new FormData(form);
   check(data);
});
*/

function GETCommand(urlCmdPath) {
   var request = new XMLHttpRequest();
   request.addEventListener('load', function(event) {
      if (request.status >= 200 && request.status < 300) {
         console.log(request.responseText);
      } else {
         console.warn(request.statusText, request.responseText);
      }
   });
   request.open("GET","http://127.0.0.1:9000/" + urlCmdPath);
   request.send();
}

function send (preset){
   GETCommand("preset/" + preset)
}

function shutdown (){
   GETCommand("exit")
}

function pureDirectOn (){
   GETCommand("pureDirect/On")
}

function pureDirectOff (){
   GETCommand("pureDirect/Off")
}

function source (sourceName){
   GETCommand("source/" + sourceName)
}