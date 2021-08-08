import json
import os, sys

class Definition:
    def __init__(self, mwresponse):
        self.raw_resonse = mwresponse

    def get_pos(self) -> list:
        pos = []
        for e in self.raw_response:
            pos.append(e.fl) 

        return list(set(pos))
    def get_def(self):
        pass

if __name__ == '__main__':
    print(f"More on {sys.argv[1]}")

