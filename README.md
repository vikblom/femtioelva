# femtioelva

## TODO

- Use different boxes for collection and rendering.
- Render as particles.
  Smoothing through simple greyscale?

# Heroku

Deploy by pushing to Heroku main.

NOTE: One gotcha was that ListenAndServe must have "0.0.0.0", else Herokus routing will not pick it up.
This is different in std HTTP and Gin.
