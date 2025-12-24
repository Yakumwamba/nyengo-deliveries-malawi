import { registerWidgetTaskHandler } from 'react-native-android-widget';
import { widgetTaskHandler } from './src/widgets';

// Register the widget task handler
// This MUST be called before AppRegistry.registerComponent
registerWidgetTaskHandler(widgetTaskHandler);
