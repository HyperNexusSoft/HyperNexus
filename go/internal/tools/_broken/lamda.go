class FireRpaMcpExtension(BaseMcpExtension):
    """Your primary task is to help users automate Android device control using AI through this MCP service."""
    route = "/firerpa/mcp/"
    name = "firerpa"
    version = "1.0"
    @mcp("tool", description="Dumps android window's layout hierarchy as JSON string.")
    def dump_window_hierarchy(self, ctx, compressed: Annotated[bool, "Enables or disables layout hierarchy compression, default true."] = True):
        data = self.device.dump_window_hierarchy(compressed).getvalue()
        return self.remove_attrs_and_empty(data)
    @mcp("tool", description="Perform a click at arbitrary coordinates on the display.")
    def click(self, ctx, pointX: Annotated[int, "X coordinate."], pointY: Annotated[int, "Y coordinate."]):
        result = self.device.click(Point(x=pointX, y=pointY))
        return str(result).lower()
    @mcp("tool", description="Perform a swipe between two points.")
    def swipe(self, ctx, fromX: Annotated[int, "Swipe-from X coordinate."], fromY: Annotated[int, "Swipe-from Y coordinate."], toX: Annotated[int, "Swipe-to X coordinate."], toY: Annotated[int, "Swipe-to Y coordinate."], step: Annotated[int, "Step to inject between two points"] = 32):
        result = self.device.swipe(Point(x=fromX, y=fromY), Point(x=toX, y=toY), step=step)
        return str(result).lower()
    @mcp("tool", description="Perform a drag between two points.")
    def drag(self, ctx, fromX: Annotated[int, "Drag-from X coordinate."], fromY: Annotated[int, "Drag-from Y coordinate."], toX: Annotated[int, "Drag-to X coordinate."], toY: Annotated[int, "Drag-to Y coordinate."]):
        result = self.device.drag(Point(x=pointX, y=pointY), Point(x=toX, y=toY))
        return str(result).lower()
    @mcp("tool", description="Get device information such as screen width, height, brand, etc.")
    def get_deviec_info(self, ctx):
        info = self.device.device_info()
        return to_json_string(MessageToDict(info))
    @mcp("tool", description="Display a toast message on the screen.")
    def show_toast(self, ctx, message: Annotated[str, "The toast message."]):
        result = self.device.show_toast(message)
        return str(result).lower()
    @mcp("tool", description="Execute script in the device's shell foreground.")
    def execute_shell_script_foreground(self, ctx, scrip: Annotated[str, "Shell script content."]):
        result = self.device.execute_script(scrip)
        return to_json_string(MessageToDict(result))
    @mcp("tool", description="Wake up the device.")
    def wake_up(self, ctx):
        result = self.device.wake_up()
        return str(result).lower()
    @mcp("tool", description="Turn off the device screen.")
    def sleep(self, ctx):
        result = self.device.sleep()
        return str(result).lower()
    @mcp("tool", description="Check if the device screen is lit up.")
    def is_screen_on(self, ctx):
        result = self.device.is_screen_on()
        return str(result).lower()
    @mcp("tool", description="Check is the device screen locked.")
    def is_screen_locked(self, ctx):
        result = self.device.is_screen_locked()
        return str(result).lower()
    @mcp("tool", description="Get the device clipboard content.")
    def get_clipboard_text(self, ctx):
        result = self.device.get_clipboard()
        return result
    @mcp("tool", description="Set the device clipboard content.")
    def set_clipboard_text(self, ctx, text: Annotated[str, "The text to set."]):
        result = self.device.set_clipboard(text)
        return str(result).lower()
    @mcp("tool", description="Simulates a short press using a key code.")
    def press_key_code(self, ctx, key_code: Annotated[int, "The Android's KeyEvent keycode."]):
        result = self.device.press_keycode(key_code)
        return str(result).lower()
    @mcp("tool", description="Get the last displayed toast on the system.")
    def get_last_toast(self, ctx):
        result = self.device.get_last_toast()
        return to_json_string(MessageToDict(result))
    @mcp("tool", description="Read android system property by name.")
    def getprop(self, ctx, name: Annotated[str, "Android system property name."]):
        return getprop(name) or ""
    @mcp("tool", description="Use full text matching to click on an element.")
    def click_by_text(self, ctx, text: Annotated[str, "The full text field."]):
        result = self.device(text=text).click_exists()
        return str(result).lower()
    @mcp("tool", code=... etc.
@mcp("tool", description="Use resourceId to input text into an input element, if the resource-id is duplicated, it cannot be used.")
def set_text_by_resource_id(self, ctx, resource_id: Annotated[str, "Input elements's resourceId (resourceName)."], text: Annotated[str, "The input text."]):
    result = self.device(resourceId=resource_id).set_text(text)
    return str(result).lower()
@mcp("tool", description="Use resourceId to input text into an input element, if the resource-id is duplicated, it cannot be used.")
def set_text_by_resource_id(self, ctx, resource_id: Annotated[str, "Input elements's resourceId (resourceName)."], text: Annotated[str, "The input text."]):
    result = self.device(res