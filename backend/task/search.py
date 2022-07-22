import glob
import os
import sys

def uncommonnamecheck(x:str):
    if "blocks" in x :
        return False
    if "account" in x :
        return False
    if "flatxs" in x :
        return True
    if "txs" in x :
        return False
    if "voteTx" in x:
        return False
    if "contractTx" in x:
        return False
    if "contracts" in x:
        return False
    

    return True


if __name__ == "__main__":
    names = [name for name in glob.glob('./src/*.tar.gz', recursive=True)]
        
    print("hoge")
    fromblock = int(sys.argv[1])
    while fromblock < int(sys.argv[2]):
        toblock = fromblock+999
        chkname = f"./src/{fromblock}_{toblock}.tar.gz"
        if chkname in names:
            pass
        else:
            print(f"Not_found_{chkname}")
        fromblock = toblock +1

    print(f"End {toblock}")