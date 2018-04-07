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

function getBasePrice(id) {
  var price = "0.00"
  switch (id) {
    case 1:
      price = "4.50"
      break;
    case 2:
      price = "2.00"
      break;
    case 3:
      price = "3.75"
      break;
    case 4:
      price = "6.50"
      break;
    default:
  }
  return price
}

function getSelectedSize() {
  if(document.getElementById("radio-short").checked) return 1
  if(document.getElementById("radio-tall").checked) return 2
  if(document.getElementById("radio-grande").checked) return 3
}

function getQuantity() {
  return parseInt(document.getElementById("quantity").value)
}

$(document).on("click", ".open-orderDialog", function () {
     var id = $(this).data('id');
     var img = $(this).data('img');
     var menu = $(this).data('menu');
     var price = getBasePrice(id)

     document.getElementById("modal-img").src = img;
     document.getElementById("modal-title").textContent = menu;
     document.getElementById("modal-title").value = id;
     document.getElementById("modal-price").textContent = price;
});

function calculatePrice() {
  var id = document.getElementById("modal-title").value
  var price = parseFloat(getBasePrice(id) * getQuantity(), 10)
  var size = getSelectedSize()
  switch (size) {
    case 1:
      price = price * 0.85
      break;
    case 2:
      break;
    case 3:
      price = price * 1.15
      break;
    default:
  }
  price = price.toFixed(2);
  document.getElementById("modal-price").textContent = price;
}

function addToCart() {
  var id = document.getElementById("modal-title").value;
  var menu = document.getElementById("modal-title").textContent;
  var img = "images/" + document.getElementById("modal-img").src.split('/').pop();
  var quantity = getQuantity()
  var size = getSelectedSize()
  var price = document.getElementById("modal-price").textContent;

  var order = '{' +
     '"id":"' + id +'",' +
     '"menu":"' + menu + '",' +
     '"img":"' + img + '",' +
     '"qty":"' + quantity + '",' +
     '"size":"' + size + '",' +
     '"price":"' + price + '"' +
  '}'
  var orders = sessionStorage.getItem('orders');
  if(orders == null) {
    var orders = '{' +
     '"orders":[' +
        order +
     ']' +
    '}'
  } else {
    orders = JSON.parse(orders);
    orders['orders'].push(JSON.parse(order));
    orders = JSON.stringify(orders);
  }
  console.log(JSON.parse(orders));
  sessionStorage.setItem('orders', orders);
  //
  // var obj = JSON.parse(order);
  // console.log(obj);
}


(function ($) {
  $('.spinner .btn:first-of-type').on('click', function() {
    $('.spinner input').val( parseInt($('.spinner input').val(), 10) + 1);
    calculatePrice()
  });
  $('.spinner .btn:last-of-type').on('click', function() {
    if (parseInt($('.spinner input').val(), 10) > 1) {
      $('.spinner input').val( parseInt($('.spinner input').val(), 10) - 1);
      calculatePrice()
    }
  });
})(jQuery);

$('body').on('hidden.bs.modal', '.modal', function () {
  $(":checkbox").prop('checked', false)
  document.getElementById("radio-tall").checked = true
  document.getElementById("quantity").value = 1
  $(this).removeData('bs.modal');
});
