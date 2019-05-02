# solrjavabindecoder
port of javabin decoder of solr in golang

This package have decoder for javabin format in solr. You may get javabin format from solr by setting wt parameter in URL to javabin.
eg: http://localhost:8983/solr/movies/select?indent=off&q=*:*&wt=javabin
Please refer to test file function BenchmarkDecodeBin to see how to use the decoder