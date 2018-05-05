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
  if(document.getElementById("radio-short").checked) return "short"
  if(document.getElementById("radio-tall").checked) return "tall"
  if(document.getElementById("radio-grande").checked) return "grande"
}

function getQuantity() {
  return parseInt(document.getElementById("quantity").value)
}

function getOptions() {
	var options = ""
	var elements = document.getElementsByTagName('input');
  for (var i = 0; i < elements.length; i++ ) {
      if (elements[i].type == 'checkbox') {
				if (elements[i].checked == true) {
          options += elements[i].value + ',';
        }
      }
  }
	return options.substring(0, options.length - 1);
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
    case "short":
      price = price * 0.85
      break;
    case "tall":
      break;
    case "grande":
      price = price * 1.15
      break;
    default:
  }
  price = price.toFixed(2);
  document.getElementById("modal-price").textContent = price;
}

function updateOrder() {
	var order = '{' +
     '"Drink":"' + menu +'",' +
     '"Size":"' + size + '",' +
     '"Options":"' + opt + '",' +
		 '"Price":' + price + ',' +
     '"Quantity":' + quantity +
  '}'
}

function deleteOrder() {
	var id = document.getElementById("modal-title").innerHTML.split(" ")[2];
	var addr = "/order/" + id
	$.ajax({
      url: addr,
      type: "DELETE",
			success: function (data) {
				window.location.replace("/status.html");
			}
  });
}


function addToCart() {
  var id = document.getElementById("modal-title").value;
  var menu = document.getElementById("modal-title").textContent;
  var img = "images/" + document.getElementById("modal-img").src.split('/').pop();
	var opt = getOptions()
  var quantity = getQuantity()
  var size = getSelectedSize()
  var price = (parseFloat(document.getElementById("modal-price").textContent) / parseFloat(quantity)).toFixed(2);
  var order = '{' +
     '"Drink":"' + menu +'",' +
     '"Size":"' + size + '",' +
     '"Options":"' + opt + '",' +
		 '"Price":' + price + ',' +
     '"Quantity":' + quantity + ',' +
     '"StoreLocation":' + sessionStorage.getItem("storeLocation")+
  '}'
  var orders = sessionStorage.getItem('orders');
  if(orders != null) {
		orders = JSON.parse(orders);
		while(quantity > 0) {
    	orders['orders'].push(JSON.parse(order));
			quantity--;
		}
    orders = JSON.stringify(orders);
		console.log(JSON.parse(orders));
		sessionStorage.setItem('orders', orders);
  }

  var addr = '/cart';
  var callType = 'post'
  if (sessionStorage.getItem("cartId") != null) {
          addr += "/"+sessionStorage.getItem("cartId");
          callType = 'put';
      }
  $.ajax({
            url: addr,
            type: callType,
            dataType: 'json',
            success: function (data) {
                sessionStorage.setItem('cartId', data.Id);
                location.reload();
            },
            data: JSON.parse(order)
        });
}

window.onload = function() {
	var orders = sessionStorage.getItem('orders');
  if(orders == null) {
    var orders = '{' +
     '"orders":[]' +
    '}'
		sessionStorage.setItem('orders', orders);
		console.log(orders);
	}
};

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
	if(!window.location.href.includes("order.html")) {
		$(":checkbox").prop('checked', false)
	}
  document.getElementById("radio-tall").checked = true
  document.getElementById("quantity").value = 1
  $(this).removeData('bs.modal');
});

window.onload
