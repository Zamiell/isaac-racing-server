$("#learn-more").click(function () {
  $("html, body").animate(
    {
      scrollTop: $("#main").offset().top,
    },
    1000
  );
});
