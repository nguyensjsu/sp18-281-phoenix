$(document).ready(function(){
	if (sessionStorage.getItem("storeLocation") != null) {
        document.getElementById("selected_store").innerHTML = "My Store: "+sessionStorage.getItem("storeLocation");
    }
});

function resizeHeaderOnScroll() {
  var distanceY = window.pageYOffset || document.documentElement.scrollTop,
  shrinkOn = 200,
  headerEl = document.getElementById('js-header');
  
  if (distanceY > shrinkOn) {
    headerEl.classList.add("smaller");
  } else {
    headerEl.classList.remove("smaller");
  }
}

window.addEventListener('scroll', resizeHeaderOnScroll);