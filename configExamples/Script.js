let tmp = execute("getProcess", {"name": "TeamViewer"}, false);
if (tmp !== "") {
   test = execute("auditpol", {"match": {"subCategoryGUID": "{0CCE9213-69AE-11D9-BED3-505054503030}"}}, false);
   if (test == "0") {
      result("demo")
   } else {
      result("failed")
   }
}

