
# python implementation with concurrent.futures.ThreadPoolExecutor

from concurrent import futures

class Solution:

    def get_domain(self, url):
        return url.split('/')[2]

    def crawl(self, startUrl: str, htmlParser: 'HtmlParser') -> List[str]:
        visited = {startUrl}
        domain = self.get_domain(startUrl)

        # using default #processors in system
        with futures.ThreadPoolExecutor() as e:
            tasks = { e.submit(htmlParser.getUrls, startUrl) }
            while tasks:
                for f in futures.as_completed(tasks):
                    tasks.remove(f)
                    for parsed_url in f.result():
                        # This is the main thread. Checking visited happen in main thread, thus don't nee locking
                        if parsed_url not in visited and self.get_domain(parsed_url) == domain:
                            visited.add(parsed_url)
                            tasks.add(e.submit(htmlParser.getUrls, parsed_url))

        return list(visited)
    