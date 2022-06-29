import glob
import os

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
    names = [name for name in glob.glob('./**/*.json', recursive=True) if uncommonnamecheck(name)]
        
    print("hoge")
    for i in names:
        if os.path.getsize(i)>512:
            print(i)