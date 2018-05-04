


function fetchDetails(){

$.get( "http://ec2-54-241-198-25.us-west-1.compute.amazonaws.com:3000/products", function( data ) {
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

