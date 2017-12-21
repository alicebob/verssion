(function($){


// dropdown
if ($(window).width() < 1200 ) {
  $('.fa-angle-down').click(function() {
    $(this).parent().find('.pure-menu-children').slideToggle();
  });
}


/* Hide Header on scroll down */

jQuery(document).ready(function( $ ){

	var didScroll;
	var lastScrollTop = 0;
	var delta = 5;
	var navbarHeight = jQuery('#header').outerHeight();

	jQuery(window).scroll(function(event){
	    didScroll = true;
	});

	setInterval(function() {
	    if (didScroll) {
	        hasScrolled();
	        didScroll = false;
	    }
	}, 250);

	function hasScrolled() {
	    var st = $(this).scrollTop();
	    
	    // Make sure they scroll more than delta
	    if(Math.abs(lastScrollTop - st) <= delta)
	        return;
	    
	    // If they scrolled down and are past the navbar, add class .nav-up.
	    // This is necessary so you never see what is "behind" the navbar.
	    if (st > lastScrollTop && st > navbarHeight){
	        // Scroll Down
	        $('#header').removeClass('nav-down').addClass('nav-up');
	    } else {
	        // Scroll Up
	        if(st + $(window).height() < $(document).height()) {
	            $('#header').removeClass('nav-up').addClass('nav-down');
	        }
	    }
	    
	    lastScrollTop = st;
	}

});

/* Menu white bg on scroll */

jQuery(window).scroll(function() {
	if ( jQuery(document).scrollTop() > 0 ) {
		jQuery('body .navbar').addClass('nav-bg');
	} else {
		jQuery('body .navbar').removeClass('nav-bg');
	}
});
    
// Match height

	$(function() {
    $('.same').matchHeight({
	    property: 'min-height'
    });
	});

/* Box slide on scroll */
	
$(window).scroll(function() {
	if ($(this).scrollTop() > 150) {
	    $('.cta.slidein').stop().animate({ right: '0px' });
	} else {
	    $('.cta.slidein').stop().animate({ right: '-390px' });
	}
});

$('.cta .close-link').click(function( event ) {
  event.preventDefault();
  $('.cta').stop().animate({ right: '-375px' }).removeClass('slidein');
});


/* Scroll to element */

	  jQuery(document).ready(function(){

		  $(".scrollto").click(function() {
          $('html, body').animate({
              scrollTop: $(".next").offset().top
          }, 1000);
          return false;
      });
	  });

// Pure Menu


/*append pure-menu-item*/
$(document).ready(function(){
  $('#menu li').addClass('pure-menu-item');
});

function resizeAgent() {
  if ($(window).width() <= 1280) {
    if($('.agent').length) {
      $('.agent').each(function(index) {
        var parent = $(this).parent();
        var child = parent.parent().find('.post:eq(1)');
        if (child.length) { //console.log($(this).height()); console.log(child.height());
          //console.log(parent.parent().find('.post:eq(1)').attr('id'));
          if ($(this).height() < child.height()) {
            $(this).css('height', child.height());
          } else {
            child.css('height', $(this).height());
          }
        }
      });
    }
  }
}

$(document).ready(function(){
  resizeAgent();
});

$(window).resize(function() {
  resizeAgent();
});


(function (window, document) {
  var menu = document.getElementById('menu-container'),
      WINDOW_CHANGE_EVENT = ('onorientationchange' in window) ? 'orientationchange':'resize';

  function toggleHorizontal() {
      [].forEach.call(
          document.getElementById('menu-container').querySelectorAll('.custom-can-transform'),
          function(el){
              el.classList.toggle('pure-menu-horizontal');
          }
      );
  };

  function toggleMenu() {
      // set timeout so that the panel has a chance to roll up
      // before the menu switches states
      if (menu.classList.contains('open')) {
          setTimeout(toggleHorizontal, 200);
      }
      else {
          toggleHorizontal();
      }
      menu.classList.toggle('open');
      document.getElementById('toggle').classList.toggle('x');
      if ($('#toggle-1').length) {
      document.getElementById('toggle-1').classList.toggle('x');
      }
  };

  function closeMenu() {
      if (menu.classList.contains('open')) {
          toggleMenu();
      }
  }

  document.getElementById('toggle').addEventListener('click', function (e) {
      toggleMenu();
  });
  if ($('#toggle-1').length) {
  document.getElementById('toggle-1').addEventListener('click', function (e) {
      toggleMenu();
  });
  }

  window.addEventListener(WINDOW_CHANGE_EVENT, closeMenu);
  })(this, this.document);

})(jQuery);
