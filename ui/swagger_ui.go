package ui

import "fmt"

func SwaggerUiHTML(config Config) string {
	return fmt.Sprintf(`
				<html lang="en">
				<head>
					<meta charset="UTF-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<meta http-equiv="X-UA-Compatible" content="ie=edge">
					<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script>
					<script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script>
					<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" />
					<title>%s</title>
					<style>
						body {
							margin: 0;
							padding: 0;
						}
					</style>
				</head>
				<body>
					<div id="swagger-ui"></div>
					<script>
						window.onload = function() {
						  SwaggerUIBundle({
							url: "%s",
							dom_id: '#swagger-ui',
							presets: [
							  SwaggerUIBundle.presets.apis,
							  SwaggerUIStandalonePreset
							],
							layout: "StandaloneLayout"
						  })
						}
					
					function addLogo() {
						const elems = document.getElementsByClassName('information-container wrapper');
						if (elems && elems.length > 0) {
							elems[0].insertAdjacentHTML("afterbegin",'<div><img src="%s" style="margin-top:24px; margin-bottom:-36px; width:200px"></div>');
							return true;
						}
						return false;
					}
					function addLogoInterval() {
						if (!addLogo()) {
							setTimeout(addLogoInterval, 10);
						}
					}
					addLogoInterval();
					</script>
				</body>
				</html>
			`, config.Title, config.SpecUrl, config.LogoUrl)
}
