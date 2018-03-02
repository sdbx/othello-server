fetch('http://localhost:8080/games/asd/actions',{
  method: 'post',
headers:{
  "Content-Type":"application/json"
}})
  .then(function(response) {
    return response.json();
  })
  .then(function(myJson) {
    console.log(myJson);
  });