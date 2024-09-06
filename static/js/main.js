jQuery(document).ready(function( $ ) {
  $('#contact-type').on('change', function() {
      var contactType = $(this).val();
      var $contactInfoInput = $('#contact-info');
      var $contactLabel = $('#contact-label');

      if (contactType === 'email') {
        $contactLabel.text('Email');
        $contactInfoInput.attr('type', 'email');
        $contactInfoInput.attr('placeholder', 'Почта для связи');
        $contactInfoInput.removeAttr('pattern');
        $contactInfoInput.removeAttr('maxlength');
      } else if (contactType === 'phone') {
        $contactLabel.text('Телефон');
        $contactInfoInput.attr('type', 'tel');
        $contactInfoInput.attr('placeholder', 'Телефон для связи');
        $contactInfoInput.attr('maxlength', '18');
        $contactInfoInput.on('input', function() {
          var value = $(this).val().replace(/\D/g, ''); // Remove all non-numeric characters
          if (value.length > 0 && !value.startsWith('7')) {
            value = '7' + value; // Ensure it starts with 7
          }
          if (value.length > 11) {
            value = value.slice(0, 11) // Limit to 11 digits
          }
          var formattedValue = '+7 ';
          if (value.length > 1) {
            formattedValue += '(' + value.slice(1, 4);
          }
          if (value.length >= 4) {
            formattedValue += ') ' + value.slice(4, 7);
          }
          if (value.length >= 7) {
            formattedValue += '-' + value.slice(7, 9);
          }
          if (value.length >= 9) {
            formattedValue += '-' + value.slice(9, 11);
          }
          $(this).val(formattedValue);
        });
      }
    });

    // Validate form on submit
  $('#contactForm').on('submit', function(e) {
    var contactType = $('#contact-type').val();
    var contactValue = $('#contact-info').val();

    if (contactType === 'email') {
      if (!validateEmail(contactValue)) {
        e.preventDefault();
        alert('Пожалуйста, введите корректный адрес электронной почты.');
      }
    } else if (contactType === 'phone') {
      if (!validatePhone(contactValue)) {
        e.preventDefault();
        alert('Пожалуйста, введите корректный номер телефона.');
      }
    }
  });

  function validateEmail(email) {
    var emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailPattern.test(email);
  }

  function validatePhone(phone) {
    var phonePattern = /^\+7\s\(\d{3}\)\s\d{3}-\d{2}-\d{2}$/;
    return phonePattern.test(phone);
  }

  $(window).scroll(function() {
    if ($(this).scrollTop() > 100) {
      $('.back-to-top').fadeIn('slow');
    } else {
      $('.back-to-top').fadeOut('slow');
    }
  });

    $('.back-to-top').click(function(){
      $('html, body').animate({scrollTop: 0}, 1500, 'easeInOutExpo');
      return false;
    });
  // Back to top button
  $(window).scroll(function() {
    if ($(this).scrollTop() > 100) {
      $('.back-to-top').fadeIn('slow');
    } else {
      $('.back-to-top').fadeOut('slow');
    }
  });
  $('.back-to-top').click(function(){
    $('html, body').animate({scrollTop : 0},1500, 'easeInOutExpo');
    return false;
  });

  $("#header").sticky({topSpacing:0, zIndex: '50'});

  $("#intro-carousel").owlCarousel({
    autoplay: true,
    dots: false,
    loop: true,
    animateOut: 'fadeOut',
    items: 1
  });

  new WOW().init();
$(document).ready(function() {
  $('.nav-menu').superfish({
    animation: {
      opacity: 'show'
    },
    speed: 400
  });

  $('.nav-menu li').removeClass('menu-active');

  var path = window.location.pathname;

  $('.nav-menu li a').each(function() {
    if ($(this).attr('href') === path) {
      $(this).parent().addClass('menu-active');
    }
  });

  if ($('#mobile-nav').length) {
    $('#mobile-nav').find('.menu-has-children').prepend('<i class="fa fa-chevron-down"></i>');
    $(document).on('click', '.menu-has-children i', function(e) {
      $(this).next().toggleClass('menu-item-active');
      $(this).nextAll('ul').eq(0).slideToggle();
      $(this).toggleClass("fa-chevron-up fa-chevron-down");
    });
  }
});
});


