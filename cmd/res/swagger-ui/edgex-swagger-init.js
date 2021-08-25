window.onload = function() {
      // Begin Swagger UI call region
      const ui = SwaggerUIBundle({
        urls: [{"url":"//localhost:8080/swagger/core-command","name":"core-command"},{"url":"//localhost:8080/swagger/core-metadata","name":"core-metadata"},{"url":"//localhost:8080/swagger/core-data","name":"core-data"},{"url":"//localhost:8080/swagger/support-notifications","name":"support-notifications"},{"url":"//localhost:8080/swagger/sys-mgmt-agent","name":"sys-mgmt-agent"},{"url":"//localhost:8080/swagger/support-scheduler","name":"support-scheduler"},{"url":"//localhost:8080/swagger/device-virtual","name":"device-virtual"},{"url":"//localhost:8080/swagger/device-rest","name":"device-rest"}],
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
      // End Swagger UI call region

      window.ui = ui;
    };