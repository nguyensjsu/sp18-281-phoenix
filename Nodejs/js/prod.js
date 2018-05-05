


function fetchDetails(){

$.get( '/products', function( data ) {
console.log(data);

 var txt = "";
for (x in data)

                  {
                    var lineStr = "Item:" + data[x].ProductName + "      " +" Category : "+ data[x].Category+" Rating : "+ data[x].Rating;
                    txt += "<li>" + lineStr + "</li>";

                  }
                  txt += "</ol>"; 
                  
document.getElementById('display').innerHTML = txt;


});
   
}

