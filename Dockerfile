FROM solr:7-alpine
ADD testDocker/movie.json /opt/
RUN start-local-solr -x && \
    bin/solr create_core -c "movies" -p 8983 && \
    bin/post -c "movies" /opt/movie.json
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["solr-foreground"]