
# python implementation using Condition, Thread

from urllib.parse import urlparse
from threading import Condition, Thread

class Solution:
    def __init__(self):
        self.queue = collections.deque()
        self.visited = set()
        self.cv = Condition()

    def can_proceed(self):
        return len(self.queue) > 0 or self.done

    def helper(self, htmlParser):
        while True:
            with self.cv:
                self.cv.wait_for(self.can_proceed)
                if self.done:
                    return
                # holding the lock 
                self.num_active += 1
                url = self.queue.popleft()

            neis = []
            for nei in htmlParser.getUrls(url):
                nei_domain = urlparse(nei).netloc 
                if nei_domain == self.domain:
                    neis.append(nei)
                        
            with self.cv:
                for nei in neis:
                    if nei not in self.visited:
                        self.visited.add(nei)
                        self.queue.append(nei)
                
                self.num_active -= 1
                if self.num_active == 0 and not self.queue:
                    self.done = True
                self.cv.notify_all()


    def crawl(self, startUrl: str, htmlParser: 'HtmlParser') -> List[str]:
        self.domain = urlparse(startUrl).netloc
        self.visited = {startUrl}
        self.queue.append(startUrl)
        
        num_thread = 10
        self.num_active = 0
        self.done = False
        threads = []
        for i in range(num_thread):
            t = Thread(target=self.helper, args=(htmlParser,))
            threads.append(t)
            t.start()

        for t in threads:
            t.join()

        return list(self.visited) 
    