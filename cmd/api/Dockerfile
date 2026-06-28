FROM hub.com/standard/golang:1.21-debian


ADD ./bin/api /home/www/github.com/neo532/gofr_layout/bin/api
ADD ./configs /home/www/github.com/neo532/gofr_layout/configs
ADD ./cmd/api/supervisord.conf /home/www/github.com/neo532/gofr_layout/supervisord.conf

WORKDIR /home/www/github.com/neo532/gofr_layout

EXPOSE 8500
EXPOSE 8501

CMD ["bin/api","-conf","configs"]
