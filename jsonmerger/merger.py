from logging import root
import pathlib
import os
import json

def main():
    working_dir = pathlib.Path(os.getcwd())
    trans_dir = os.path.join(working_dir,"src/assets/i18n")

    files_info = read_files(trans_dir)
    dest_dir = os.path.join(working_dir,"output")
    pathlib.Path(dest_dir).mkdir(exist_ok=True)
    
    for pkg in files_info['packages']:
        pkgjson = {}
        for file in files_info['files']:
            if file['package'] == pkg:
                key = file['key']
                pkgjson[key] = file['content']
        
        out_name = 'exp.'+pkg+'.json'
        out_file = os.path.join(dest_dir,out_name)
        
        with open(out_file, "w") as outfile:
            content = json.dumps(pkgjson, indent=4, ensure_ascii=False)
            outfile.write(content)

def read_files(rootdir:str)->dict :
    result = {}
    files_info = []
    packages = {}
    for path, subdirs, files in os.walk(rootdir):
        for name in files:
            filepath = pathlib.Path(os.path.join(path, name))
            if filepath.suffix == ".json":
                finfo = {
                    "name":name,
                    "path":str(filepath),
                    "package":filepath.parent.name,
                    "parent": str(filepath.parent),
                    "ext": filepath.suffix,
                    "key": name.replace(filepath.suffix,"")
                }
                packages[finfo['package']]=True
                with open(finfo["path"],'r') as json_file:
                    finfo["content"] = json.load(json_file)
                
                files_info.append(finfo)
    
    result['files'] = files_info
    result["packages"] = packages

    return result

if __name__ == "__main__":
    main()