/* @license
 * Copyright (C) Symbolic Software — All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * Written by Nadim Kobeissi <nadim@symbolic.software>
 */

#import "MobileEbitenViewControllerWithErrorHandling.h"

#import <Foundation/Foundation.h>

@implementation MobileEbitenViewControllerWithErrorHandling {
}

- (void)onErrorOnGameUpdate:(NSError*)err {
    // You can define your own error handling e.g., using Crashlytics.
    NSLog(@"Inovation Error!: %@", err);
}

@end
