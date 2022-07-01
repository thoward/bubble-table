# Dynamic Filter Example

This example shows how to use a dynamically generated footer. This allows full control of the table footer to show content that changes with the table.

In this example, the table implements [vertical scrolling](../vertical_scrolling) and also a [custom filter input](../filterapi) located outside of the table. Normally this would mean that the default footer would duplicate the filter information, or would need to be static content.

With the dynamic footer, we show that the footer can change it's contents depending on the status of the filter (how many rows shown vs in the dataset). Using this method, any kind of changing information can be shown in the footer.
