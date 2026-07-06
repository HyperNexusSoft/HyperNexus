from pathlib import Path
import re

content = Path("apps/web/src/app/dashboard/DashboardHomeClient.tsx").read_text()
# Ensure React Fragment is correctly implemented
if 'activeTab === \'browser\' && (' in content:
    print("Browser tab conditional is present")
